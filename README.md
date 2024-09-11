# dbschema-anygen

dbschema-anygen is useful go generator that mainly uses schema of database.

You can pass any templates to the generator.

## How to use

Implement the code below into your generator code.

```golang
cfg := api.Config{
    // You should set configuration values
    ...
}

generator := api.NewGenerator(api.WithFuncMap(
    map[string]interface{}{
        "YourFunction": YourFunction,
    },
))
err := generator.Generate(c.Context, cfg)
```
