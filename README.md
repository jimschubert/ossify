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

## Examples

### Licenses

#### Find licenses by keyword

License keywords

* copyleft
* discouraged
* international
* miscellaneous
* non-reusable
* obsolete
* osi-approved
* permissive
* popular
* redundant
* retired
* special-purpose

Example:

```shell script
ossify license --keyword popular
Apache-2.0          (Apache License, Version 2.0)
BSD-2               (BSD 2-Clause License)
BSD-3               (BSD 3-Clause License)
CDDL-1.0            (Common Development and Distribution License, Version 1.0)
EPL-1.0             (Eclipse Public License, Version 1.0)
GPL-2.0             (GNU General Public License, Version 2.0)
GPL-3.0             (GNU General Public License, Version 3.0)
LGPL-2.1            (GNU Lesser General Public License, Version 2.1)
LGPL-3.0            (GNU Lesser General Public License, Version 3.0)
MIT                 (MIT/Expat License)
MPL-2.0             (Mozilla Public License, Version 2.0)
```

#### Write out license text to file

```shell script
ossify license MIT > LICENSE
# or, with --id switch. Both are case-insensitive.
ossify license --id mit > LICENSE
```

#### Search for specific license text

```shell script
ossify license --search apache
Apache-1.1          (Apache Software License, Version 1.1)
Apache-2.0          (Apache License, Version 2.0)
```

#### View details for specific license

```shell script
ossify license MIT --details
MIT                 (MIT/Expat License)
osi-approved, popular, permissive

Common names
  * MIT
  * Expat

License Standards
  * DEP5       MIT
  * DEP5       Expat
  * SPDX       MIT
  * Trove      License :: OSI Approved :: MIT License

Links
  * https://opensource.org/licenses/mit (OSI Page)
  * https://tldrlegal.com/license/mit-license (tl;dr legal)
  * https://en.wikipedia.org/wiki/MIT_License (Wikipedia page)
```

## License

This project is [Licensed MIT](./LICENSE)

All included license text is [Licensed CC0 1.0 Universal](./data/licenses/LICENSE.CC0)
