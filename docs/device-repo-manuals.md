# Device repo manuals folder structure

When the source folder does not have a specific manual-type for a device, but there is a device folder with a `details.json` file, the repository_url will be used to find a 'manufacturer' version of that manual in the repository.

## Providing manuals in a device firmware repository

Anyone that makes firmware for a device, can supply manuals for that device firmware.

This can be done by using the following folder structure from the root of the repository.

### Folder structure

```text
.
└── docs/
    └── manuals/
        └── <manual_type>/
            ├── languages/
            │   ├── nl-NL.md
            │   └── en-GB.md
            └── assets/
                └── ...
```

### `devices` directory

The main directory where the device-specific manuals are placed.

### `manual-type` directory

A folder for a manual type (e.g. `info`, or `installation`). 

It contains a `languages` folder and possibly an `assets` folder.

#### `languages/`

This folder contains markdown files for different languages. The names of the markdown files must be: `<language_code>.md` (e.g. `nl-NL.md` or `en-GB.md`).

#### `assets/`

This folder contains assets that can be reffered to from manuals made in the `languages` folder. All assets should be placed here, to make sure they are handled correctly.