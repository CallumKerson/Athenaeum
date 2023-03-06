/**
 EditDateView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct EditDateView: View {
    @ObservedObject var viewModel: EditDateViewModel

    init(viewModel: inout EditDateViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        DatePicker(selection: $viewModel.dateState, in: ...Date(), displayedComponents: .date) {
            Text(viewModel.label)
        }
        .onChange(of: viewModel.dateState) { _ in
            viewModel.wasEdited()
        }
    }
}

struct EditDateView_Previews: PreviewProvider {
    static var releaseDateViewModel: EditDateViewModel = .init(
        label: "Release Date",
        initialDate: ISO8601DateFormatter().date(from: "2021-04-25T08:00:00Z")
    )

    static var previews: some View {
        EditDateView(viewModel: &releaseDateViewModel)
    }
}
