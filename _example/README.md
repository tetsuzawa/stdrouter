# Example of Standard Library Router Generator

1. Write routing config to `router.go`
2. Execute command as below.
   ```shell script
    $ stdrouter
    ```
3. `router_gen.go` will be created. This is the implementation of router.
   
# Tips

`NewRouter` returns `http.Handler`. So you can use middleware func.






