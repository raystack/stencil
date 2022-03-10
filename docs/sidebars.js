module.exports = {
  docsSidebar: [
    'introduction',
    'usecases',
    'installation',
    'glossary',
    {
      type: "category",
      label: "Guides",
      items: [
        "guides/introduction",
        "guides/quickstart",
        "guides/manage_namespace",
        "guides/manage_schemas",
        "guides/clients",
      ],
    },
    {
      type: "category",
      label: "Formats",
      items: [
        "formats/protobuf",
        "formats/avro",
        "formats/json",
      ],
    },
    {
      type: "category",
      label: "Server",
      items: [
        "server/overview",
        "server/rules",
      ],
    },
    {
      type: "category",
      label: "Clients",
      items: [
        "clients/overview",
        "clients/go",
        "clients/java",
        "clients/clojure",
        "clients/js",
      ],
    },
    {
      type: "category",
      label: "Reference",
      items: [
        "reference/api",
        "reference/cli",
      ],
    },
    {
      type: "category",
      label: "Contribute",
      items: [
        "contribute/contribution",
      ],
    },
    'roadmap',
  ],
};