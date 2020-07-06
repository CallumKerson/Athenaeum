/**
 PublicationDate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

struct PublicationDate: Equatable, Codable, Hashable {
    public let year: Int
    public let month: Int?
    public let day: Int?

    public var getDateAsString: String? {
        if let month = month {
            if let day = day {
                return "\(year)-\(String(format: "%02d", month))-\(String(format: "%02d", day))"
            } else {
                return "\(year)-\(String(format: "%02d", month))"
            }
        }
        return String(year)
    }

    public var asDate: Date? {
        var components = DateComponents()
        components.year = year
        components.month = month ?? 1
        components.day = day ?? 1
        components.hour = 8
        components.timeZone = TimeZone(identifier: "UTC")!
        return Calendar(identifier: Calendar.Identifier.iso8601)
            .date(from: components)
    }

    public init(year: Int, month: Int?, day: Int?) {
        self.year = year
        self.month = month
        self.day = day
    }

    public init(from dateString: String) throws {
        let detector = try? NSDataDetector(types: NSTextCheckingResult.CheckingType.date.rawValue)
        let matches = detector?.matches(in: dateString, range: NSMakeRange(0, dateString.utf16.count))
        let date = matches?.first?.date
        if let date = date {
            let calendar = Calendar.current
            year = calendar.component(.year, from: date)
            month = calendar.component(.month, from: date)
            day = calendar.component(.day, from: date)
            return
        }
        throw DateParseError.notAValidDate(dateString)
    }
}

enum DateParseError: Error {
    case notAValidDate(_ date: String)
}
