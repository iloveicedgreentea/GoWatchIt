Config package

It has functions to read and write configs to the DB using reflect

This allows you to provide any struct and it will be able to process it

structs need db tags

helper functions exist to provide simple access to values

globalConfig is a global var in memory so configs can be access without passing around a ton of stuff to every function