/**
 EditDateViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

class EditDateViewModel: ObservableObject {
    @Published var dateState: Date
    let label: String
    var hasBeenEdited: Bool
    var hadExistingValue: Bool

    init(label: String, initialDate: Date?) {
        self.label = label
        if let initialDate = initialDate {
            self.dateState = initialDate
            self.hadExistingValue = true
        } else {
            self.hadExistingValue = false
            self.dateState = Date()
        }
        self.hasBeenEdited = false
    }

    func date() -> Date? {
        if self.hadExistingValue || self.hasBeenEdited {
            logger.info("Date is being saved")
            return self.dateState
        } else {
            logger.info("Date is not being saved")
            return nil
        }
    }

    func wasEdited() {
        self.hasBeenEdited = true
    }
}
