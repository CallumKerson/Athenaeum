/**
 Double+StringRepresentation.swift
 Copyright (c) 2020 Callum Kerr-Edwards
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
