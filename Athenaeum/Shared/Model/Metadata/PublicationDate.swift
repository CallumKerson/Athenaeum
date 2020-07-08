/**
 PublicationDate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct PublicationDate: Equatable, Codable, Hashable {
    public let year: Int
    public let month: Int?
    public let day: Int?

    public var getDateAsString: String? {
        if let month = month {
            if let day = day {
                return "\(self.year)-\(String(format: "%02d", month))-\(String(format: "%02d", day))"
            } else {
                return "\(self.year)-\(String(format: "%02d", month))"
            }
        }
        return String(self.year)
    }

    public var asDate: Date? {
        var components = DateComponents()
        components.year = self.year
        components.month = self.month ?? 1
        components.day = self.day ?? 1
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
        let matches = detector?.matches(
            in: dateString,
            range: NSMakeRange(0, dateString.utf16.count)
        )
        let date = matches?.first?.date
        if let date = date {
            let calendar = Calendar.current
            self.year = calendar.component(.year, from: date)
            self.month = calendar.component(.month, from: date)
            self.day = calendar.component(.day, from: date)
            return
        }
        throw DateParseError.notAValidDate(dateString)
    }
}

enum DateParseError: Error {
    case notAValidDate(_ date: String)
}
