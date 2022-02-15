# Compatability rules

Stencil server provides following compatibility rules.
- BACKWARD_COMPATABILITY
- FORWARD_COMPATABILITY
- FULL_COMPATABILITY

Stencil currently supports protobuf, avro and json schema formats. Compatibility rules for each schema format has been built separately considering each schema format's features.

## Feature support matrix

| Compatability rule | Protobuf | Avro | JSON |
| ------------------ | --- | --- | --- |
| BACKWARD_COMPATABILITY | Yes | Yes | No |
| FORWARD_COMPATABILITY | Yes | Yes | No |
| FULL_COMPATABILITY | Yes | Yes | No |

## Protobuf compatibility rules

Protobuf compatability rules composed of [compatability checks](#list-of-checks).

### Rules

| Compatibility name | List of checks |
| ------------------ | -------------- |
| BACKWARD_COMPATABILITY | SYNTAX_CHANGE, MESSAGE_DELETE, NON_INCLUSIVE_RESERVED_RANGE, NON_INCLUSIVE_RESERVED_NAMES, FIELD_DELETE, FIELD_JSON_NAME_CHANGE, FIELD_LABEL_CHANGE, FIELD_KIND_CHANGE, FIELD_TYPE_CHANGE, ENUM_DELETE, ENUM_VALUE_DELETE, ENUM_VALUE_NUMBER_CHANGE |
| FORWARD_COMPATABILITY | SYNTAX_CHANGE, MESSAGE_DELETE, NON_INCLUSIVE_RESERVED_RANGE, NON_INCLUSIVE_RESERVED_NAMES, FIELD_JSON_NAME_CHANGE, FIELD_LABEL_CHANGE, FIELD_KIND_CHANGE, FIELD_TYPE_CHANGE, FIELD_DELETE_WITHOUT_RESERVED_NUMBER, FIELD_DELETE_WITHOUT_RESERVED_NAME, ENUM_DELETE, ENUM_VALUE_NUMBER_CHANGE, ENUM_VALUE_DELETE_WITHOUT_RESERVEDNUMBER, ENUM_VALUE_DELETE_WITHOUT_RESERVEDNAME |
| FULL_COMPATABILITY | SYNTAX_CHANGE, MESSAGE_DELETE, NON_INCLUSIVE_RESERVED_RANGE, NON_INCLUSIVE_RESERVED_NAMES, FIELD_DELETE, FIELD_JSON_NAME_CHANGE, FIELD_LABEL_CHANGE, FIELD_KIND_CHANGE, FIELD_TYPE_CHANGE, ENUM_DELETE, ENUM_VALUE_DELETE, ENUM_VALUE_NUMBER_CHANGE |

### List of Checks

| Check | Description |
| ---- | ------------ |
| SYNTAX_CHANGE | checks if proto file syntax does not switch between proto2 and proto3, including going to/from unset (which assumes proto2) to set to proto3. Changing the syntax results in differences in generated code for many languages. |
| MESSAGE_DELETE | checks that messages are deleted from a given file. Deleting a message will delete the corresponding generated type, which could be referenced in source code. Instead of deleting these types, deprecate them using [`deprecated` option](https://developers.google.com/protocol-buffers/docs/proto3#options).|
| NON_INCLUSIVE_RESERVED_NAMES | Checks if current reserved names contains all previous reserved names. |
| NON_INCLUSIVE_RESERVED_RANGE | Checks if current reserve range inclusive of previous reserved range. This check ensures previous reserved tag numbers haven't been deleted |
| FIELD_DELETE | checks that no message field is deleted. Deleting message field will result in the field being deleted from the generated source code, which could be referenced. Instead of deleting these, deprecate them using [`deprecated` option](https://developers.google.com/protocol-buffers/docs/proto3#options). |
| FIELD_JSON_NAME_CHANGE | Checks if the json_name for field does not change, which would break JSON compatibility. |
| FIELD_LABEL_CHANGE | checks that no field changes it's label, i.e. `optional`, `required`, `repeated`. Changing to/from optional/required and repeated will be a generated source code and JSON breaking change. Changing to/from optional and repeated is actually not a wire-breaking change, however changing to/from optional and required is. Given that it's unlikely to be advisable in any situation to change your label, and that there is only one exception, we find it best to just outlaw this entirely. |
| FIELD_KIND_CHANGE | checks that a field has the same type. Changing the type of a field can affect the type in the generated source code, wire compatibility, and JSON compatibility. |
| FIELD_TYPE_CHANGE |  Checks if message/enum field it's message/enum type has changed from previous version. This rule only applies to message kind and enum kind. |
| FIELD_DELETE_WITHOUT_RESERVED_NUMBER | Checks if field is deleted, it's tag number should be added to reserved numbers. This will ensure deleted field tag number won't be used in future. |
| FIELD_DELETE_WITHOUT_RESERVED_NAME | Checks if field is deleted, it's tag name should be added to reserved names. This will help to keep the JSON compatibility. |
| ENUM_DELETE | Checks that no enum is deleted from current version. Deleting an enum will delete the corresponding generated type, which could be referenced in source code. Instead of deleting these types, deprecate them using `deprecated` option. |
| ENUM_VALUE_DELETE | Checks that no enum value is deleted. Deleting an enum value will result in the corresponding value being deleted from the generated source code, which could be referenced. Instead of deleting these, deprecate them. |
| ENUM_VALUE_DELETE_WITHOUT_RESERVEDNUMBER | Checks if enum value deleted, it's enum number should be added to reserved numbers. |
| ENUM_VALUE_DELETE_WITHOUT_RESERVEDNAME | Checks if enum value deleted, it's enum name should be added to reserved names. This will help to keep the JSON compatability |
| ENUM_VALUE_NUMBER_CHANGE | Check if enum number has changed between current, previous versions. For example You cannot change FOO_ONE = 1 to FOO_ONE = 2. Doing so will result in potential JSON incompatibilites and broken source code. |






