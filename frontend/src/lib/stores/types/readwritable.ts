import {writable} from "svelte/store";

export default function readWritable<T>(initial: T) {
    const {subscribe, set, update} = writable(initial);

    let cached = initial;
    return {
        subscribe,
        update: (fn: (v: T) => T) => {
            update(v => {
                const retVal = fn(v);
                cached = retVal;
                return retVal;
            });
        },
        set: (val: T) => {
            cached = val;
            set(val);
        },
        get value() {return cached;}
    };
};
