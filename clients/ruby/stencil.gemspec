# frozen_string_literal: true

require_relative "lib/stencil/version"

Gem::Specification.new do |spec|
  spec.name          = "stencil"
  spec.version       = Stencil::VERSION
  spec.authors       = ["Daval Pargal"]
  spec.email         = ["davalpargal@gmail.com"]

  spec.summary       = "Stencil ruby gem provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date."
  spec.homepage      = "https://odpf.gitbook.io/stencil/"
  spec.required_ruby_version = ">= 2.4.0"

  spec.metadata["allowed_push_host"] = ""

  spec.metadata["homepage_uri"] = spec.homepage
  spec.metadata["source_code_uri"] = "https://github.com/odpf/stencil"
  spec.metadata["changelog_uri"] = "https://github.com/odpf/stencil/blob/master/CHANGELOG.md"

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`.split("\x0").reject { |f| f.match(%r{\A(?:test|spec|features)/}) }
  end
  spec.bindir        = "exe"
  spec.executables   = spec.files.grep(%r{\Aexe/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.17.3"
  spec.add_development_dependency "webmock", "~> 3.14.0"
end
