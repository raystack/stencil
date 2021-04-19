package io.odpf.stencil.utils;

public class RandomUtils {

    public long getRandomNumberInRange(long min, long max) {

        if (min >= max) {
            throw new IllegalArgumentException("max must be greater than min");
        }

        return (long) Math.random()*((max - min) + 1L) + min;
    }
}
