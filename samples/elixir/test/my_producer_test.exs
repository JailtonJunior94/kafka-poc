defmodule MyProducerTest do
  use ExUnit.Case
  doctest MyProducer

  test "greets the world" do
    assert MyProducer.hello() == :world
  end
end
