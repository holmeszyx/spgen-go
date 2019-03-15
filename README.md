Using toml to define configurations, And generating source code with it.


### Usage

create new configuration file

```
spgen -new config.toml
```

generate source code

```
spgen -o std config.toml
spgen -o android:kt config.toml
```