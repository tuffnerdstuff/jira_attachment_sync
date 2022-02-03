# Jira Attachment Sync
Jira Attachment Sync (`jas`) is a simple tool for downloading and extracting all attachments associated wih a specific Jira issue. 

## Configuration
`jas` requires some basic configuration in order to work properly. The path to a TOML file containing the configuration can be passed by using the option `--configPath "/path/to/config.toml"`. If `--configPath` is not set it defaults to `"./jas-config.toml"`. You can use `jas-config_template.toml` as a template.

## Usage
Once `jas` is configured you can execute it by providing a jira issue key by passing the option `--issue`.

## Example
Download attachments of `BUG-123` and use config at `/path/to/config.toml`:
`jas --configDir /path/to/config.toml --isue BUG-123`