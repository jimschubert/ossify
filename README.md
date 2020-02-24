# ossify

![build](https://github.com/jimschubert/ossify/workflows/build/badge.svg) ![goreleaser](https://github.com/jimschubert/ossify/workflows/goreleaser/badge.svg)

<blockquote>
<dl>
<dt><em>ossify (n.)</em></dt>
<dd>
    <ul>
    <li><s>to change (a material, such as cartilage) into bone</s></li>
    <li><strong>to make rigidly conventional and opposed to change</strong></li>
    </ul>
</dd>
</dl>
</blockquote>

---

:warning: This is an early project that I'm using to become familiar with Go. Contributions are welcome, but this isn't yet to be considered "ready" for anyone to even look at.

Open Source Software often contain common project layouts and components across multiple projects, languages, or groups.
This tool aims to provide a means to automate and/or verify projects using conventions; those built into the tool and those defined by the user.

See [KriaSoft/Folder-Structure-Conventions](https://github.com/KriaSoft/Folder-Structure-Conventions) and [Standard Go Project Layout](https://github.com/golang-standards/project-layout) for some examples of conventional directory structure.

This project aims to go one step further and include templating and configuration for common files in open source software. Examples include:

* `.gitignore`
* `README.md`
* `LICENSE`
* `CONTRIBUTING.md`
* GitHub Issues and Pull Request templates

## License

This project is [Licensed MIT](./LICENSE)

All included license text is [Licensed CC0 1.0 Universal](./data/licenses/LICENSE.CC0)
