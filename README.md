# xsync
Implementations of extra sync package modifications

1. Once.Do  - https://github.com/golang/go/issues/22098

This is for cases where initialization function may fail due to external reasons. For an exactly once case, if the function fails
due to external reason, exponential backoff also doesn't help as there is still a case for failure, and after that
the execution will not trigger the initialization again 

This implementation tries to guarrantee exactly one **successful** execution of the initialization function. 
