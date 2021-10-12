require 'simplecov'

SimpleCov.formatter = SimpleCov::Formatter::MultiFormatter.new(
  [
    SimpleCov::Formatter::HTMLFormatter,
  ]
)

SimpleCov.start do
  add_filter "/spec/"
  minimum_coverage 80
end

require 'bundler/setup'
Bundler.setup

require 'pry'
require 'stencil'

RSpec.configure do |config|
  config.filter_run_when_matching focus: true
  config.disable_monkey_patching!
  config.profile_examples = 10
  config.order = :random
end

module JSON
  module_function

  def reader(file_name)
    path = File.join(File.dirname(__FILE__), "test_data/#{file_name}")
    File.read(path)
  end
end
