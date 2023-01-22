## How to Use

Using aquila is a 3-step process:

### Setup Config Files

- Create `aquila.yaml` config file.
- Specify the `codePath` & `docsPath` directory in the config file.

### Setting Up

- Annotate all the code regions that might be used within the docs using a comment `#aqstart <region-name>`
  and closing it with `#aqend <region-name>`.
- Wherever you want to use the code region in the docs, use the following
  syntax: `<?code-region <file-name> region="<region-name>"?>`.

### Using

- Generate the docs using `aquila generate` command.
- The generated docs will be present in the `./code_regions` directory.
- Run `aquila update-md` to update the docs with the generated code regions.

Example can be found [here](example/).

