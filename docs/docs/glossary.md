# Glossary

This section describes the core elements of a schema registry.

## Namespace

A named collection of schemas. Each namespace holds a logically related set of schemas, typically managed by a single entity, belonging to a particular application and/or having a shared access control management scope. Since a schema registry is often a resource with a scope greater than a single application and might even span multiple organizations, it is very useful to put a grouping construct around sets of schemas that are related either by ownership or by a shared subject matter context. A namespace has following attributes:

- **ID:** Identifies the schema group.
- **Format:** Defines the schema format managed by this namespace. e..g Avro, Protobuf, JSON
- **Compatibility** Schema compatibility constraint type. e.g. Backward, Forward, Full

## Schema

A document describing the structure, names, and types of some structured data payload. Conceptually, a schema is a description of a data structure. Since data structures evolve over time, the schema describing them will also evolve over time. Therefore, a schema often has multiple versions.

## Version

A specific version of a schema document. Even though not prescribed in this specification, an implementation might choose to impose compatibility constraints on versions following the initial version of a schema.

## Compatibility

A key Schema Registry feature is the ability to version schemas as they evolve. Compatibility policies are created at the namespace or schema level, and define evolution rules for each schema.

After a compatibility policy has been defined for a schema, any subsequent version updates must honor the schema’s original compatibility, to allow for consistent schema evolution.

Compatibility of schemas can be configured with any of the below values:

### Backward

Indicates that new version of a schema would be compatible with earlier version of that schema.

### Forward

Indicates that an existing schema is compatible with subsequent versions of the schema.

### Full

Indicates that a new version of the schema provides both backward and forward compatibilities.
