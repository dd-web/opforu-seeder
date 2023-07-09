# OPforu Seeder
Performance optimized database seeding utility. Seeds 25k+ cross referenced/linked data structures in under a second, using a batch & reference pointer methods to generate huge amounts of useful data very fast.

>**WARNING**: this will _**drop**_ any database it's given access to in order to generate it's data. Please use caution

## About
This is the result of an abstraction of a more general purpose useful data generation library. The aim of the tool is to be a minimal setup database seeding utility with an emphasis on performance. Input a schema, set your references and constraints, then let the seeder do the rest. Creating a tool as described is more of an undertaking than I could've imagined, so this serves as the proof of concept and a demonstration of sorts.

Currently uses MongoDB but the tool I am working on supports Postgres and MySQL engines as well.


## Usage
pull down the repo and build the binary.

```bash
make build
```

or if you want to run it immediately

```bash
make run
```

Ensure you have the correct .env variables set to establish a connection to mongodb.
