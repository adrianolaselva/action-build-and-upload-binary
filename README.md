## "Upload Application Binary" Action for GitHub Actions

Upload application in Github.

## Usage

```yaml
name: example
on: [push]

jobs:
  hello_world_job:
    runs-on: ubuntu-latest
    name: A job to say hello
    steps:
      - name: Hello world action step
        uses: action-build-and-upload-binary@master
        env:
          API_TOKEN_GITHUB: ${{ secrets.API_TOKEN_GITHUB }} 
          OWNER: adrianolaselva
          REPOSITORY: repository-name
          TAG: '0.0.1'
```