# NeedForHeat manual server
Server that server NeedForHeat manuals written in Markdown.

## Table of contents
* [Deploying](#deploying)
* [Developing](#developing)
* [Usage](#usage)
* [Features](#features)
* [Status](#status)
* [License](#license)
* [Citation](#citation)
* [Credits](#credits)

## Deploying
For our process to deploy the API to our public server, or update it, see these links:
- Deploy: https://github.com/energietransitie/twomes-backoffice-configuration#manuals
- Update: https://github.com/energietransitie/twomes-backoffice-configuration#updating

### Prerequisites
The NeedForHeat manual server is available as a Docker image.
You will need to [install Docker](https://docs.docker.com/engine/install/) to run it.

### Images
See all [available images](https://github.com/energietransitie/needforheat-manual-server/pkgs/container/needforheat-manual-server):
- Use the `latest` tag to get the latest stable release built from a tagged GitHub release. 
- Use the `main` tag to get the latest development release, built directly from the `main` branch.

### Docker Compose ([more information](https://docs.docker.com/compose/features-uses/))
```yaml
version: "3.8"
services:
  web:
    container_name: needforheat-manual-server
    image: ghcr.io/energietransitie/needforheat-manual-server:latest
    ports:
      - 8080:8080
    volumes:
      - ./source:/source
    environment:
      - NFH_MANUAL_SOURCE=https://github.com/energietransitie/needforheat-manuals.git
      - NFH_MANUAL_SOURCE_BRANCH=tst
      - NFH_FALLBACK_LANG=en-US
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

### Writing manuals

Manuals are written in Markdown. The Markdown files should be placed in a folder structure following some rules. Read [this](./docs/source-folder-structure.md) document to see what the folder structure has to be, which files to place in it and what names to give them.

Manuals can be written by device firmware makers. Read [this](./docs/device-repo-manuals.md) document to see how you can write manuals for a specific device when making firmware for it.

### Device display names
Device display names can be retrieved from `/devices/<device-name>`.

This will return a json object with display names in different languages.

### Device manuals
Device manuals can be retrieved from `/devices/<device-name>/<manual-type>`.

A generic manual (not specific to a campaign) can be retrieved from `/devices/<device-name>/`. This will auto redirect to `/devices/<device-name>/generic`.

If you try to retrieve a manual of a specific type that does not exist, you will be redirected to the 'manufacturer' version if it exists.

Manuals from `/devices/<device-name>/<manual-type>` will automatically redirect to the language that your browser requests using the Accept-Language header. e.g. `/devices/<device-name>/<manual-type>/en-US/` for a British English version.

### Campaign manuals
Campaign manuals can be retrieved from `/campaigns/<campaign-name>/<manual-type>`.

A generic manual (not specific to a campaign, but more generic to the lab e.g. privacy policy) can be retrieved from `/campaigns/<manual-type>/`. This will auto redirect to `/campaigns/generic/<manual-type>`.

Manuals from `/campaigns/<campaign-name>/<manual-type>` will automatically redirect to the language that your browser requests using the Accept-Language header. e.g. `/campaigns/<campaign-name>/<manual-type>/en-US/` for a British English version.

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

## Citation

If you use this repository in your research or work, please cite the following pre-print, which describes the overall NeedForHeat DataGear system of which this repository is a part:

> Ter Hofte, H., & van Ravenzwaaij, N. (2025). *NeedForHeat DataGear: An Open Monitoring System to Accelerate the Residential Heating Transition*. arXiv preprint arXiv:2509.06927. https://doi.org/10.48550/arXiv.2509.06927

**Note:** This is a pre-print submitted on 8 Sep 2025 and has not yet been peer-reviewed. For the full paper, see [https://arxiv.org/abs/2509.06927](https://arxiv.org/abs/2509.06927).

## Credits
This software is created by:
* Nick van Ravenzwaaij · [@n-vr](https://github.com/n-vr)

Product owner:
* Henri ter Hofte · [@henriterhofte](https://github.com/henriterhofte)

Thanks also goes to:
* Harris Mesic - [@Labhatorian](https://github.com/Labhatorian)

We use and gratefully acknowlegde the efforts of the makers of the following source code and libraries:
* [chi](https://github.com/go-chi/chi), by Peter Kieltyka, Google Inc, licensed under [MIT license](https://github.com/go-chi/chi/blob/master/LICENSE)
* [markdown](https://github.com/gomarkdown/markdown), by Russ Ross, Krzysztof Kowalczyk, licensed under [BSD 2-clause license](https://github.com/gomarkdown/markdown/blob/master/LICENSE.txt)
* [go-git](https://github.com/go-git/go-git), by Sourced Technologies, S.L., lincensed under [Apache 2.0 license](https://github.com/go-git/go-git/blob/master/LICENSE)
