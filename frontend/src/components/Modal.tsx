import React from "react";

type Props = {
  onClose: () => void;
  onSubmit: (name: string, lat: number, lon: number) => void;
  mode: "add" | "edit";
  initialData?: { name: string; lat: number; lon: number };
};

export default function LocationModal({
  onClose,
  onSubmit,
  mode,
  initialData,
}: Props) {
  const [name, setName] = React.useState(initialData?.name || "");
  const [lat, setLat] = React.useState(initialData?.lat.toString() || "");
  const [lon, setLon] = React.useState(initialData?.lon.toString() || "");

  // ...handleClick sama

  const handleSubmit = () => {
    if (name && lat && lon) {
      onSubmit(name, parseFloat(lat), parseFloat(lon));
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
      <div
        id="modal-box"
        className="bg-white p-6 rounded-lg w-[300px] space-y-4"
      >
        <h2 className="text-lg font-bold">
          {mode === "add" ? "Add Location" : "Edit Location"}
        </h2>
        <input
          className="w-full border px-3 py-1 rounded"
          placeholder="City Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input
          className="w-full border px-3 py-1 rounded"
          placeholder="Latitude"
          value={lat}
          onChange={(e) => setLat(e.target.value)}
        />
        <input
          className="w-full border px-3 py-1 rounded"
          placeholder="Longitude"
          value={lon}
          onChange={(e) => setLon(e.target.value)}
        />
        <div className="flex justify-end gap-2">
          <button className="px-3 py-1 bg-gray-300 rounded" onClick={onClose}>
            Cancel
          </button>
          <button
            className="px-3 py-1 bg-blue-500 text-white rounded"
            onClick={handleSubmit}
          >
            {mode === "add" ? "Submit" : "Edit"}
          </button>
        </div>
      </div>
    </div>
  );
}
