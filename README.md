# go-gitlab
Gitlab api client for golang applications

## Installing
```
go get github.com/kryabinin/go-gitlab
```

## Example
Get user information
```go
client := gitlab.NewClient("token")

user, err := client.GetUserByID(context.Background(), 1234)
if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
}

fmt.Println("user's name is: ", user.Name)
```

