# OPforu Seeder

This project has changed goals and now remains as a single solution seeder to a specific data set.

>**WARNING**: this will _**drop**_ any database it's given access to in order to generate it's data. Please use caution

## About

This was originally intended as a proof of concept which worked quite well. The idea was to input a schema and have a generative data seeder. There came many complexities and ultimately I came to the conclusion it's far easier just to do this again and adapt it to any usecase.

The main idea is to eliminate the network call wherever possible, so the entirety of the data is generated and calls bulk store operations. This tool only cares about the final state of the data so if you're using migrations you'll need to know the final state of the schema. It keeps track of references in a shared store, so they may be updated/retreived as necessary. It also means if you're going to use this approach it should be designed carefully to ensure references are always valid when accessed. 

## Usage

pull down the repo. build the binary:

```bash
make build
```

to run:

```bash
make run
```

Ensure you have the correct .env variables set to establish a connection to the database.


## Notes

Yeah so this is just a seeder for an old project I used to learn web frameworks. I'm leaving it up here because I hate seeding for 20 minutes and figure some other people might want to do something similar for themselves. 

If you want to take this and adapt it to your own project feel free to do so.
