/**
 EditSeriesView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import SwiftUI

struct EditSeriesView: View {
    @ObservedObject var viewModel: EditSeriesViewModel

    init(viewModel: inout EditSeriesViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        HStack {
            Text("Series")
            Spacer()
            TextField("", text: $viewModel.seriesTitleState)
            Spacer()
            Text("Entry Number")
            TextField("", text: $viewModel.seriesEntryState)
                .onReceive(Just(viewModel.seriesEntryState)) { newValue in
                    let filtered = newValue.filter { "0123456789.".contains($0) }
                    if filtered != newValue {
                        self.viewModel.seriesEntryState = filtered
                    }
                }
        }
    }
}

struct EditSeriesView_Previews: PreviewProvider {
    static var emptySeriesViewModel: EditSeriesViewModel = .init(
        initialTitle: nil,
        initalEntry: nil
    )

    static var existingSeriesViewModel: EditSeriesViewModel = .init(
        initialTitle: "The Stormlight Archive",
        initalEntry: 2.5
    )

    static var previews: some View {
        Group {
            EditSeriesView(viewModel: &emptySeriesViewModel)
            EditSeriesView(viewModel: &existingSeriesViewModel)
        }
    }
}
