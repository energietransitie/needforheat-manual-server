# Source folder structure

This describes the folder structure that you need to use in order to create manuals.

## General guidelines

## Categories

Manuals can fall into two categories, each with its own folder structure:
- Device manuals
- Campaign manuals

## Device manuals

Device manuals are manuals for a specific device. This can be something like an installation manual or FAQ. They are placed in the `devices` directory of the manual source.

The manuals can be different per campaign. To customise a manual for a specific campaign, place it in a directory with the name of that campaign. To make a manual that is not for a specific campaign, place it in a folder with the name `generic`.

### Folder structure

```text
.
└── devices/
    └── <device-type>/
        ├── details.json
        ├── display_names.json
        └── <manual-type>/
            └── <campaign-name>/
                ├── languages/
                │   ├── nl-NL.md
                │   └── en-US.md
                └── assets/
                    └── ...
```

### `devices` directory

The main directory where the device-specific manuals are placed.

### `device-type` directory

A folder with the name of a device. The folder contains types of manuals (e.g. `installtion`-manual, or `info`-manual).

This folder also contains two files: `details.json` and `display_names.json`.

#### `details.json`

This file has the firmware repository that is used to get the 'manufacturer' manuals that are privided in the repo.

```json
{
    "firmware_repository": "https://github.com/org/repo"
}
```

#### `display_names.json`

This file has the human readable display name for each supported language.

```json
{
  "nl-NL": "Slimme meter module",
  "en-US": "Smart meter module"
}
```

### `manual-type` directory

A folder for a manual type (e.g. `installtion`-manual, or `info`-manual). This folder contains folders for specific campaigns, or `generic` (no campaign).

### `campaign-name` directory

This folder contains a version of a manual that is specific for a campaign (or not, if the folder name is `generic`).

It contains a `languages` folder and possibly an `assets` folder.

> Note that the name `manufacturer` is reserved and should not be used.

#### `languages/`

This folder contains markdown files for different languages. The names of the markdown files must be: `<language_code>.md` (e.g. `nl-NL.md` or `en-US.md`).

#### `assets/`

This folder contains assets that can be reffered to from manuals made in the `languages` folder. All assets should be placed here, to make sure they are handled correctly.

## EnergyQuery manuals

EnergyQueries follow the same structure as devices.

## Campaign manuals

Campaign manuals are manuals for a specific campaign. This can be something like a privacy policy or an FAQ. They are placed in the `campaigns` directory of the manual source.

### Folder structure

```text
.
└── campaigns/
    └── <campaign-name>/
        └── <manual-type>/
            ├── languages/
            │   ├── nl-NL.md
            │   └── en-US.md
            └── assets/
                └── ...
```

### `campaigns` directory

The main directory where the campaign manuals are placed.

### `campaign-name` directory

This folder contains manuals for a specific campaign (or not, if the folder name is `generic`).

### `manual-type` directory

A folder for a manual type (e.g. `privacy-policy`, or `faq`). 

It contains a `languages` folder and possibly an `assets` folder.

#### `languages/`

This folder contains markdown files for different languages. The names of the markdown files must be: `<language_code>.md` (e.g. `nl-NL.md` or `en-US.md`).

#### `assets/`

This folder contains assets that can be reffered to from manuals made in the `languages` folder. All assets should be placed here, to make sure they are handled correctly.
