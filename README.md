# What is carmensandiego?
`carmensandiego` (Secret scanner named after our favorite detective, writen in sandie-*Go*) is a **set of tools and scripts** that aims to provide you with a continuous application security posture management toolset based on a **constantly generated**, **lightweight**, **quickly actionable** set of vulnerability reports, saved onto a file you can process after the scanning is done.

# Why do we need this?
Secret scanning is an essential part of the application security management task list, and that requires a proper toolset. There are open tools that provide this out there like `Trivy`, or repositories with large datasets of [secret pattern patching rules](https://github.com/mazen160/secrets-patterns-db), but these are either hard to scale in a considerable codebase, or need a large amount of work to manage.

# Why create this tool if Trivy exists?
Trivy is very good at its job of investigating commited secrets in a repository against a set of standards, like SSH keys or authentication tokens, we can target and clear. That said, Trivy is designed to be launched manually, is not optimized for larger codebases, and its results checked in terminal interfaces. It's really useful to have a security scanner that we can run to check for secrets in bulk, but a 500 line `json` pasted in an ephimeral terminal after is not the **human-centric** way of working. This is where `carmensandiego` comes in.

# What does atomic-threat-hunter bring to the table?
The aim of this tool is to bring in the human-centric part we miss from other tool sources, like `Trivy` or `secrets-patterns-db:`:
- **Ease of use**: Configure the tool once with the profiles you need, and leave it running forever to work for you
- **A dataset you can work with**: send the data you got, anywhere. Standard format, simple processing.
- **Narrow down on your target**: Is the general scope too heavy on your target? make a faster, more lightweight profile to run simultaneously. Scale as you need.
- **Speed as a center point**: Avoiding spinning rust and sequential execution. It's not 1999 anymore. You got RAM, use it. You got cores, use them.
- **Portable and reproducible**: docker image, ready to build and test!

# How does this work?
This CLI tool is built entirely as a Go binary, with all its dependencies built inside a Docker image with the provided `Dockerfile`. The image uses the `helper` built binary as an entrypoint, so just build it:
```
docker build -t carmensandiego
```

modify your desired profile in `config.json`, and run the container with the desired task and profile:
```
docker run --rm \
    -v ./config.json:/carmensandiego/config.json \
    -v ./scans:/carmensandiego/scans \
    carmensandiego:latest scan example-profile
```

# Where can I send the detected findings?
The standard format provided by `carmensandiego` can be sent to anything that then supports the data processing a standard `json` file.