# go-ftd

Go Bindings for [Cisco FirePower NGFW](https://www.cisco.com/c/en/us/products/collateral/security/firepower-ngfw/datasheet-c78-736661.html). These bindings talk to Firepower Device Manager.

## Example

Open a Session using env vars:

```go
params := make(map[string]string)
params["grant_type"] = "password"
params["username"] = os.Getenv("FTD_USER")
params["password"] = os.Getenv("FTD_PASSWORD")
params["debug"] = "true"
params["insecure"] = "true"

ftd, err := NewFTD(os.Getenv("FTD_HOST"), params)
if err != nil {
    glog.Errorf("error: %s\n", err)
    return nil, err
}

return ftd, nil
```

Creating a Network Object:

```go
// Create a Network Object for a single host 1.1.1.1
n := new(NetworkObject)
n.Name = "testObj001"
n.SubType = "HOST"
n.Value = "1.1.1.1"

err = ftd.CreateNetworkObject(n, DuplicateActionReplace)
if err != nil {
    glog.Errorf("error: %s\n", err)
    return
}
```

Creating an Access Rule:

```go
// Allow any traffic between any and network object n1 and network object group g1
a := new(AccessRule)
a.Name = "testPolicy001"
a.RuleAction = RuleActionPermit
a.EventLogAction = LogActionNone
// n1.Refence() returns a reference object of a Network Object
a.DestinationNetworks = append(a.DestinationNetworks, n1.Reference())
// g1.Refence() returns a reference object of a Network Object Group
a.DestinationNetworks = append(a.DestinationNetworks, g1.Reference())

err = ftd.CreateAccessRule(a, "default")
if err != nil {
    glog.Errorf("error: %s\n", err)
    return
}
```

## Authors

* **Remi Philippe** - *Initial work* - [remiphilippe](https://github.com/remiphilippe)

See also the list of [contributors](https://github.com/remiphilippe/go-ftd/contributors) who participated in this project.

## License

This project is licensed under the Apache 2 License - see the [LICENSE](LICENSE) file for details