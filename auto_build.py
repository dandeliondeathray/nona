#!/usr/bin/env python3

"""Automatic builds of all containers."""

import subprocess
import sys


def container_exists(name):
    """Check if a container already exists."""
    return subprocess.call(["docker", "inspect", name],
                           stdout=subprocess.DEVNULL,
                           stderr=subprocess.DEVNULL) == 0


def build_image(name, directory, build_args):
    """Build a Docker image."""
    builds_args_cmdline = [["--build-arg", "{}".format(arg)] for arg in build_args]
    command_line = ["docker", "build", "-t", name, directory] + sum(builds_args_cmdline, [])
    print("Building Docker container with:", ' '.join(command_line))
    return subprocess.call(command_line,
                           stdout=None, #subprocess.DEVNULL,
                           stderr=None) # subprocess.DEVNULL)


class BuildError(Exception):
    """Raised when a build step fails."""
    def __init__(self, builder, context):
        message = "FAILED: {}, commit_id = {}".format(builder, context.commit_id)
        Exception.__init__(self, message)


class Builder:
    """Base class for anything that builds stuff."""
    def __init__(self, dependencies):
        self.dependencies = dependencies

    def build(self, context):
        """Build all dependencies for this builder."""
        for dep in self.dependencies:
            if not dep.exists(context):
                dep.build(context)
            else:
                print("Artifact for {} already exists".format(dep))


class Context:
    """Context for building stuff."""
    def __init__(self, commit_id):
        self.commit_id = commit_id

    def __str__(self):
        return "Context({})".format(self.commit_id)

class DockerPusher(Builder):
    """Push Docker images."""
    def __init__(self, images_to_push, dependencies=None):
        Builder.__init__(self, dependencies or None)
        self._images_to_push = images_to_push

    def exists(self, _context):
        """DockerPusher produces no artifacts, so run it every time. It's idempotent."""
        return False

    def build(self, context):
        """Push images to Docker Hub."""
        Builder.build(self, context)
        for image_name in self._images_to_push:
            contextual_image_name = image_name.format(commit_id=context.commit_id)
            exit_code = subprocess.call(["docker", "push", contextual_image_name])
            if exit_code != 0:
                raise BuildError(self, context)

    def __str__(self):
        return "Push to Docker Hub:" + '\n\t'.join(self._images_to_push)


class DockerTagger(Builder):
    """Tag an image."""
    def __init__(self, from_name, to_name, dependencies=None):
        Builder.__init__(self, dependencies or [])
        self._from_name = from_name
        self._to_name = to_name

    def exists(self, _context):
        """Tagging is idempotent anyway."""
        return False

    def build(self, context):
        """Tag an image."""
        Builder.build(self, context)
        from_name_in_context = self._from_name.format(commit_id=context.commit_id)
        to_name_in_context = self._to_name.format(commit_id=context.commit_id)
        exit_code = subprocess.call(["docker", "tag", from_name_in_context, to_name_in_context])
        if exit_code != 0:
            raise BuildError(self, context)

    def __str__(self):
        return "Tag {} as {}".format(self._from_name, self._to_name)


class DockerImageBuilder(Builder):
    """Builds a Docker image."""
    def __init__(self, name, directory, build_args, dependencies=None):
        Builder.__init__(self, dependencies or [])
        self._image_name = name
        self._directory = directory
        self._build_args = build_args

    def exists(self, context):
        """Check if a built artifact already exists."""
        return container_exists(self._image_name.format(commit_id=context.commit_id))

    def build(self, context):
        """Build this Docker image."""
        Builder.build(self, context)
        real_build_args = [arg.format(commit_id=context.commit_id) for arg in self._build_args]
        image_name = self._image_name.format(commit_id=context.commit_id)
        exit_code = build_image(image_name,
                                self._directory,
                                real_build_args)
        if exit_code != 0:
            raise BuildError(self, context)

    def __str__(self):
        return "Build Docker Image: {} in {} with args {}".format(self._image_name,
                                                                  self._directory,
                                                                  self._build_args)


class OnlyDependencies(Builder):
    """Build nothing, but collect dependencies into a single builder."""
    def __init__(self, dependencies):
        Builder.__init__(self, dependencies)

    def exists(self, _context):
        """Nothing to create, so can't exist."""
        return False

    def build(self, context):
        """Only build dependencies"""
        Builder.build(self, context)

    def __str__(self):
        return "OnlyDependencies()"


PLUMBER = DockerImageBuilder("erikedin/nonaplumber:{commit_id}",
                             "service/plumber",
                             [])
NONAINTERFACE = DockerImageBuilder("erikedin/nonainterface:{commit_id}",
                                   "service/nonainterface",
                                   [])
PUZZLESTORE = DockerImageBuilder("erikedin/nonapuzzlestore:{commit_id}",
                                 "service/puzzlestore",
                                 ["PLUMBER_TAG={commit_id}"],
                                 dependencies=[PLUMBER])
SLACKMESSAGING = DockerImageBuilder("erikedin/nonaslackmessaging:{commit_id}",
                                    "service/slackmessaging",
                                    ["PLUMBER_TAG={commit_id}"],
                                    dependencies=[PLUMBER])
NONACONTROL = DockerImageBuilder("erikedin/nonacontrol:{commit_id}",
                                 "service/control",
                                 ["INTERFACE_TAG={commit_id}"],
                                 dependencies=[NONAINTERFACE])
TAG_PUZZLESTORE_AS_STAGING = DockerTagger("erikedin/nonapuzzlestore:{commit_id}",
                                          "erikedin/nonapuzzlestore:staging",
                                          dependencies=[PUZZLESTORE])
TAG_SLACKMESSAGING_AS_STAGING = DockerTagger("erikedin/nonaslackmessaging:{commit_id}",
                                             "erikedin/nonaslackmessaging:staging",
                                             dependencies=[SLACKMESSAGING])
TAG_NONACONTROL_AS_STAGING = DockerTagger("erikedin/nonacontrol:{commit_id}",
                                          "erikedin/nonacontrol:staging",
                                          dependencies=[NONACONTROL])

TAG_AS_STAGING = OnlyDependencies(dependencies=[TAG_PUZZLESTORE_AS_STAGING,
                                                TAG_SLACKMESSAGING_AS_STAGING,
                                                TAG_NONACONTROL_AS_STAGING])
PUSH_TO_STAGING = DockerPusher(["erikedin/nonapuzzlestore:staging",
                                "erikedin/nonaslackmessaging:staging",
                                "erikedin/nonacontrol:staging"],
                               dependencies=[TAG_AS_STAGING])


def main():
    """Entry to point build all containers."""
    try:
        commit_id = subprocess.check_output(["git", "rev-parse", "HEAD"]).decode('UTF-8')
        commit_id = commit_id.strip()
    except subprocess.CalledProcessError:
        print("Could not get commit id. Aborting.")
        sys.exit(1)
    context = Context(commit_id)
    print("Building with context", context)
    try:
        PUSH_TO_STAGING.build(context)
    except BuildError as build_error:
        print(build_error)

if __name__ == "__main__":
    main()
