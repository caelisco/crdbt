Archived project due to license changes in CockroachDB and no longer needing a utility to work with the database.
# crdbt
CockroachDB Tools (crdbt) is a simple command line utility that helps manage the installation, configuration and upgrading of CockroachDB.

Intergreatme makes use of CockroachDB. However, it is often tedious to remember all the commands required to interact with CockroachDB when trying to upgrade from one version to another.

I had already written several scripts to help manage CockroachDB as well as help assist with the upgrade process.

I opted to combine the scripts in to a single Go application.

**Important:**
Early versions of crdbt are aimed at only satisfying my own needs to interact and manage CockroachDB at Intergreatme.
It is:
- Designed to run on Linux (Ubuntu x64)
- Works with systemd
