/**
 Double+StringRepresentation.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

extension Double {
    var asString: String {
        if truncatingRemainder(dividingBy: 1.0) == 0.0 {
            return "\(Int(self))"
        } else {
            return "\(self)"
        }
    }
}
