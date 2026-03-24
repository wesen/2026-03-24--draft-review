import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface ReaderState {
  started: boolean;
  currentSectionId: string | null;
  readSectionIds: string[];
}

const initialState: ReaderState = {
  started: false,
  currentSectionId: null,
  readSectionIds: [],
};

export const readerSlice = createSlice({
  name: "reader",
  initialState,
  reducers: {
    startReading(state, action: PayloadAction<string>) {
      state.started = true;
      state.currentSectionId = action.payload;
    },
    goToSection(state, action: PayloadAction<string>) {
      state.currentSectionId = action.payload;
    },
    markSectionRead(state, action: PayloadAction<string>) {
      if (!state.readSectionIds.includes(action.payload)) {
        state.readSectionIds.push(action.payload);
      }
    },
    resetReader() {
      return initialState;
    },
  },
});

export const { startReading, goToSection, markSectionRead, resetReader } =
  readerSlice.actions;
