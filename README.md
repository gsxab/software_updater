# Software Update Watcher

---

A crawler to fetch updates of pieces of software,
written for software which is not used frequently enough to check updates or unable to,
for where a software package manager or a marketplace is unavailable,
and for those who want to manage updates themselves.

---

## Usage

### Knowledge Prerequisites

- YAML and JSON.
- Database operations.

### Configuration

A YAML in the working directory is required to provide program configurations.
`conf.yaml.example` is a template.

### Database Preparation

The editor is still on the way, so direct operations on database is unavoidable by now...

Only sqlite3 is supported by now.
Run the gorm/gen subproject in `core/db/main` to generate a new database file.

For any software, fields `name`, `homepage_url`, and `version_actions` are required.

### Front End

If you run the project from source, please remember to compile the front-end vue project,
and copy the distribution directory into files under a `dist` folder,
and compile it together to make a bundled application.

### Running

The software is composed of a server (this repo) and a front-end project.
Access the main page configured.
And the front-end would list all the pieces of software and make basic operation.



---

## License

The software is distributed under the GNU GPLv3 License,
or any later version.

Note that abuse of crawlers may bring trouble to site maintainers,
and may cause legal issues in some countries and regions.
We expect any users to obey their local laws;
any illegal, immoral usage of the software is ABSOLUTELY NOT SUPPOSED. 

---

## For Developer

Some features are still working on now or to do in the future:
- Not all actions are tested now.
- Flow editor, at least a copy-paste-tune tool.
- Plugin for RPC to common download tools, e.g. aria2.
Some features are not planned by now, for restricted manpower, which seems helpful in using.
- Home page information editor.
- Plugin mechanics for the version page and the flow page.
- More useful pages?
Some features are blocked now.
- Multi-branch flows. Hard to visualize and edit.
- Optional and repetition flow nodes. Not useful enough in single-branch flows.
