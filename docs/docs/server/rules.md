# Compatability rules

Stencil server provides following compatibility rules

## BACKWARD_COMPATABILITY

This rule checks for the following conditions,
- **File delete**: no file has been deleted. Deleting a file will result in it's generated header file being deleted as well, which could break source code.
- **Package change**: file has the same package value. Changing the package value will result in a lot of issues downstream in various languages
- **Syntax change**: file does not switch between proto2 and proto3, including going to/from unset (which assumes proto2) to set to proto3. Changing the syntax results in differences in generated code for many languages.
- **Options change**: no change in java package, JAVA outer class name and go package file options.


## MESSAGE_NO_DELETE

This rule checks that messages are deleted from a given file. Deleting a message will delete the corresponding generated type, which could be referenced in source code. Instead of deleting these types, deprecate them.

For example,
```
message Foo {
  option deprecated = true;
}
```

## FIELD_NO_BREAKING_CHANGE

This rule checks for the following conditions.
- **Field delete**: checks that no message field is deleted. Deleting message field will result in the field being deleted from the generated source code, which could be referenced. Instead of deleting these, deprecate them.
for example,
  ```
  message Bar {
    string one = 1 [deprecated = true];
  }
  ```
- **Number change**: checks that field number for specific field is not changed. For example, you cannot change int64 foo = 1; to int64 bar = 1;. This affects generated source code, but also affects JSON compatibility as JSON uses field names for serialization. This does not affect wire compatibility, however we generally don't recommend changing field names.
- **Type change**: checks that a field has the same type. Changing the type of a field can affect the type in the generated source code, wire compatibility, and JSON compatibility.
- **Label change**: This checks that no field changes it's label, i.e. `optional`, `required`, `repeated`. Changing to/from optional/required and repeated will be a generated source code and JSON breaking change. Changing to/from optional and repeated is actually not a wire-breaking change, however changing to/from optional and required is. Given that it's unlikely to be advisable in any situation to change your label, and that there is only one exception, we find it best to just outlaw this entirely.
- **JSON name change**: This checks that the json_name field option does not change, which would break JSON compatibility.


## ENUM_NO_BREAKING_CHANGE

This rule checks for the following conditions.
- **Enum delete**: Checks that no enum is deleted from current version. Deleting an enum will delete the corresponding generated type, which could be referenced in source code. Instead of deleting these types, deprecate them. For example,
  ```
  enum Foo {
    option deprecated = true;
    FOO_UNSPECIFIED = 0;
    ...
  }

  ```
- **Enum value delete**: This checks that no enum value is deleted. Deleting an enum value will result in the corresponding value being deleted from the generated source code, which could be referenced. Instead of deleting these, deprecate them.
  ```
  enum Foo {
    FOO_UNSPECIFIED = 0;
    FOO_ONE = 1 [deprecated = true];
  }

  ```
- **Enum value number change**: This checks that a given enum value has the same name for each enum value number. For example You cannot change FOO_ONE = 1 to FOO_TWO = 1. Doing so will result in potential JSON incompatibilites and broken source code.






