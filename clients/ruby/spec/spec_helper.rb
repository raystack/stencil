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
require 'webmock/rspec'

RSpec.configure do |config|
  config.filter_run_when_matching focus: true
  config.disable_monkey_patching!
  config.profile_examples = 10
  config.order = :random
end

WebMock.disable_net_connect!(allow_localhost: true)