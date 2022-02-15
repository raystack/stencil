module Stencil
  RSpec.describe Store do
    context "#read" do
      let(:sample_key) { "sample_key" }
      let(:sample_value) { 123 }

      it "should be able to handle concurrent reads of data" do
        store = Store.new
        store.write(sample_key, sample_value)

        5.times do
          Thread.start do
            expect(store.read(sample_key)).to eq(sample_value)
          end
        end
      end
    end

    context "#write" do
      let(:sample_key) { "sample_key" }
      let(:sample_value) { 100 }

      it "should be able to handle concurrent writes of data by locking data" do
        store = Store.new
        store.write(sample_key, sample_value)
        expect(store.read(sample_key)).to eq(sample_value)
      end
    end
  end
end
