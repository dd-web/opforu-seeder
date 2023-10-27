# OPforu Seeder

This project has changed goals and now remains as a single solution seeder to a specific data set.

>**WARNING**: this will _**drop**_ any database it's given access to in order to generate it's data. Please use caution

## About

This was originally intended as a proof of concept which worked quite well, squeezing around 4million records per sec on a single thread. The idea was to input a schema and have a generative data seeder. There came many problems and ultimately I came to the conclusion it's far easier just to do this again and adapt to any usecase.

if you want to copy this go ahead. The main idea is to eliminate the network calls. because the GO compiler is smart enough to inline the function chain the performance is pretty signifigant on a single thread and memory usage is low even though it's purely a utility that _holds data_ until it fits the schema.

I haven't tested it but i'm 99% sure you're going to lose performance multithreading something like this. if using SQL you can generate a schema from migrations (without inputting any data yet) down to a single schema state and generate your data from there. Asyncronous blocks even sequential operations. if you don't NEED it, don't use it.

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

I understand the program needs to keep good track of the entire database state but it's lifetime is small. it runs for a second and it's gone. Right now it runs a bit longer than that because it's running an encryption on each password for each account, and it's an expensive function. other than that it's fast.

If you want to take this and adapt it to your own project feel free to do so.