/**
 EditSeriesViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

class EditSeriesViewModel: ObservableObject {
    @Published var seriesTitleState: String
    @Published var seriesEntryState: String

    init(initialTitle: String?, initalEntry: Double?) {
        self.seriesTitleState = initialTitle ?? ""
        self.seriesEntryState = (initalEntry ?? 0).clean
    }

    func series() -> Series? {
        guard let entry = Double(seriesEntryState) else { return nil }
        if self.seriesTitleState.isEmpty || entry == 0 {
            logger.info("Series is empty")
            return nil
        } else {
            logger.info("Series is not empty")
            return Series(entry: entry, title: self.seriesTitleState)
        }
    }
}

extension Double {
    var clean: String {
        self.truncatingRemainder(dividingBy: 1) == 0 ? String(format: "%.0f", self) : String(self)
    }
}
