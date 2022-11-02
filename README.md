# crdbt
crdbt is a command line utility to help work with Cockroach DB

Cockroach DB Tools (crdbt) is a simple command line utility that helps manage the installation, configuration and upgrading of Cockroach DB.

At Intergreatme, we make extensive use of Cockroach DB. However, it is often tedious to remember all the commands required to interact with Cockroach DB when trying to upgrade from one version to another.

I had already written several scripts to help manage Cockroach DB as well as help assist with the upgrade process.

I opted to combine the scripts in to a single Go application.

** Important: 
early versions of crdbt are aimed at only satisfying my own needs to interact and manage Cockroach DB at Intergreatme. It is designed to run on Linux, and has only been tested to run on Ubuntu, and works with systemd.