# Twomes manual server
Server that server twomes manuals written in Markdown.

## Table of contents
* [Deploying](#deploying)
* [Developing](#developing)
* [Usage](#usage)
* [Features](#features)
* [Status](#status)
* [License](#license)
* [Credits](#credits)

## Deploying
For our process to deploy the API to our public server, or update it, see these links:
- Deploy: https://github.com/energietransitie/twomes-backoffice-configuration#manuals
- Update: https://github.com/energietransitie/twomes-backoffice-configuration#updating

### Prerequisites
The Twomes manual server is available as a Docker image.
You will need to [install Docker](https://docs.docker.com/engine/install/) to run it.

### Images
See all [available images](https://github.com/energietransitie/twomes-manual-server/pkgs/container/twomes-manual-server):
- Use the `latest` tag to get the latest stable release built from a tagged GitHub release. 
- Use the `main` tag to get the latest development release, built directly from the `main` branch.

### Docker Compose ([more information](https://docs.docker.com/compose/features-uses/))
```yaml
version: "3.8"
services:
  web:
    container_name: twomes-manual-server
    build: .
    ports:
      - 8080:8080
    volumes:
      - ./source:/source
    environment:
      - TWOMES_MANUAL_SOURCE=./source
      - TWOMES_FALLBACK_LANG=en-GB
```

## Developing
This section describes how you can change the source code using a development environment and compile the source code into a binary release of the firmware that can be deployed, either via the development environment, or via the method described in the section [Deploying](#deploying).

### Prerequisites
- [Go (minimum 1.20)](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

### Running
Make sure Docker is running on your local machine, then start the service from a command line terminal from the root of this repository:
```shell
docker compose up --build
```

This generates log messages.
Just keep this running in your terminal.

The manuals are now available on http://localhost:8080/.

## Usage

### Device manuals
Device manuals can be retrieved from `/devices/<device-name>/<manual-type>`.

A generic manual (not specific to a campaign) can be retrieved from `/devices/<device-name>/`. This will auto redirect to `/devices/<device-name>/generic`.

If you try to retrieve a manual of a specific type that does not exist, you will be redirected to the 'manufacturer' version if it exists.

Manuals from `/devices/<device-name>/<manual-type>` will automatically redirect to the language that your browser requests using the Accept-Language header. e.g. `/devices/<device-name>/<manual-type>/en-GB/` for a British English version.

### Campaign manuals
Campaign manuals can be retrieved from `/campaigns/<campaign-name>/<manual-type>`.

A generic manual (not specific to a campaign, but more generic to the lab e.g. privacy policy) can be retrieved from `/campaigns/<manual-type>/`. This will auto redirect to `/campaigns/generic/<manual-type>`.

Manuals from `/campaigns/<campaign-name>/<manual-type>` will automatically redirect to the language that your browser requests using the Accept-Language header. e.g. `/campaigns/<campaign-name>/<manual-type>/en-GB/` for a British English version.

## Features
Ready:
* Parse markdown files to HTML manuals.
* Serve HTML files.
* Redirect to correct language based on Accept-Language header.
* Redirect to generic campaign if none is specified.
* Redirect to device firmware repository for missing manuals.
* Manual source can be set to local directory or git repository.

To-do:
* Watching for file changes and update without restarting the server.
* Get page titles from display_names.json for language automatically when generating HTML.
* A friendly "manual not found" (404) page that can contain contact information if desired.
* Support authentication for private git repositories.

## Status
Project is: _in progress_

## License
This software is available under the [Apache 2.0 license](./LICENSE), Copyright 2023 [Research group Energy Transition, Windesheim University of Applied Sciences](https://windesheim.nl/energietransitie) 

## Credits
This software is created by:
* Nick van Ravenzwaaij · [@n-vr](https://github.com/n-vr)

Product owner:
* Henri ter Hofte · [@henriterhofte](https://github.com/henriterhofte) · Twitter [@HeNRGi](https://twitter.com/HeNRGi)

We use and gratefully acknowlegde the efforts of the makers of the following source code and libraries:
* [chi](https://github.com/go-chi/chi), by Peter Kieltyka, Google Inc, licensed under [MIT license](https://github.com/go-chi/chi/blob/master/LICENSE)
* [markdown](https://github.com/gomarkdown/markdown), by Russ Ross, Krzysztof Kowalczyk, licensed under [BSD 2-clause license](https://github.com/gomarkdown/markdown/blob/master/LICENSE.txt)
* [go-git](https://github.com/go-git/go-git), by Sourced Technologies, S.L., lincensed under [Apache 2.0 license](https://github.com/go-git/go-git/blob/master/LICENSE)
