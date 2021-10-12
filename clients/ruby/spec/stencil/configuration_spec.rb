module Stencil
  RSpec.describe Configuration do
    let(:configuration) { Stencil.configuration }

    it 'should set registry urls correctly' do
      expect(configuration.registry_urls).to eq([])

      expected_value = ['http://localhost:3000']
      configuration.registry_urls = expected_value
      expect(configuration.registry_urls).to eq(expected_value)
    end


    it 'should set http_timeout correctly' do
      expect(configuration.http_timeout).to eq(DEFAULT_TIMEOUT_IN_MS)

      timeout = 3000
      configuration.http_timeout = timeout
      expect(configuration.http_timeout).to eq(3000)
    end

    it 'should set refresh_enabled correctly' do
      expect(configuration.refresh_enabled).to eq(true)

      refresh_enabled = false
      configuration.refresh_enabled = refresh_enabled
      expect(configuration.refresh_enabled).to eq(refresh_enabled)
    end

    it 'should set refresh_ttl_in_secs correctly' do
      expect(configuration.refresh_ttl_in_secs).to eq(DEFAULT_REFRESH_INTERVAL_IN_SECONDS)

      refresh_ttl_in_secs = 50000
      configuration.refresh_ttl_in_secs = refresh_ttl_in_secs
      expect(configuration.refresh_ttl_in_secs).to eq(refresh_ttl_in_secs)
    end

    it 'should set bearer_token correctly' do
      token = "sampletoken"
      expected_bearer_token = "Bearer sampletoken"
      configuration.bearer_token = token
      expect(configuration.bearer_token).to eq(expected_bearer_token)
    end



    describe '#configure' do
      let(:refresh_enabled) {true}
      let(:refresh_ttl_in_secs) {60000}
      let(:registry_urls) {["abc.com/latest"]}
      let(:token) {"ABCD1234"}
      let(:bearer_token) {"Bearer " + token}
      let(:http_timeout) {6000}

      before(:each) do
        Stencil.configure do |config|
          config.registry_urls = registry_urls
          config.bearer_token = token
          config.refresh_enabled = refresh_enabled
          config.refresh_ttl_in_secs = refresh_ttl_in_secs
          config.http_timeout = http_timeout
        end
      end

      subject { Stencil.configuration }

      it 'should set configuration correctly' do
        expect(subject.registry_urls).to eq(registry_urls)
        expect(subject.bearer_token).to eq(bearer_token)
        expect(subject.refresh_enabled).to eq(refresh_enabled)
        expect(subject.refresh_ttl_in_secs).to eq(refresh_ttl_in_secs)
        expect(subject.http_timeout).to eq(http_timeout)
      end
    end
  end
end
