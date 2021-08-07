# predicate
Add predicates to github issues and automatically close issues if a predicate is met.


## Use
- Install this action on your repo with schedule trigger
- When an issue is created with a `predicate` codeblock the issue will be closed when the predicate exits with a 0 exit code:

### Issue example
```
This is a description


    ```predicate
    curl https://google.com
    ```
```


![img.png](img.png)