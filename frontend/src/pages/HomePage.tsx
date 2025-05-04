import { useState, useEffect } from "react";
import { Plus, Pencil, BarChart2, Trash2Icon } from "lucide-react";
import LocationModal from "../components/Modal";
import axios from "axios";
import toast from "react-hot-toast";
import {
  getWeatherIconUrl,
  getWeatherDescription,
  getWeatherCodeFromSummary,
} from "../utils/WeatherIcons";

export default function HomePage() {
  const [locations, setLocations] = useState<any[]>([]);
  const [selected, setSelected] = useState<number[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [editingLocation, setEditingLocation] = useState<any | null>(null);
  const [flippedCard, setFlippedCard] = useState<number | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await axios.get("http://localhost:8080/locations", {
          headers: { Authorization: `Bearer ${token}` },
          withCredentials: true,
        });

        const locationsRaw = res.data;

        // Ambil predictions untuk semua lokasi
        const locationsWithPredictions = await Promise.all(
          locationsRaw.map(async (item: any) => {
            const summary = item.weather_summary?.Valid
              ? item.weather_summary.String
              : "Not available";

            const predictionRes = await axios.get(
              `http://localhost:8080/locations/${item.id}/predictions`,
              {
                headers: { Authorization: `Bearer ${token}` },
                withCredentials: true,
              }
            );

            return {
              id: item.id,
              name: item.name,
              lat: item.latitude,
              lon: item.longitude,
              summary: summary,
              code: getWeatherCodeFromSummary(summary),
              temperature: item.temperature?.Valid
                ? item.temperature.Float64
                : null,
              windspeed: item.wind_speed?.Valid
                ? item.wind_speed.Float64
                : null,
              updatedAt: item.updated_at,
              predictions: predictionRes.data || [], // ⬅️ tambah prediksi
            };
          })
        );

        setLocations(locationsWithPredictions);
      } catch (err) {
        console.error("Failed to fetch locations or predictions:", err);
      }
    };

    fetchData();
  }, []);

  const toggleSelect = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const toggleFlip = (id: number) => {
    setFlippedCard((prev) => (prev === id ? null : id));
  };

  const handleDelete = async () => {
    try {
      const token = localStorage.getItem("token");

      for (const id of selected) {
        await axios.delete(`http://localhost:8080/locations/${id}`, {
          headers: { Authorization: `Bearer ${token}` },
          withCredentials: true,
        });
      }

      toast.success("Berhasil menghapus lokasi.");
      setSelected([]);
      setTimeout(() => window.location.reload(), 2000);
    } catch (error) {
      console.error("Gagal menghapus lokasi.", error);
      toast.error("Gagal menghapus lokasi.");
    }
  };

  return (
    <div className="relative flex flex-col items-center h-screen bg-weather">
      <h1 className="text-[50px] font-extrabold pt-2 flex">WEATHER APP</h1>

      <button
        onClick={() => setShowModal(true)}
        className="bg-amber-500 hover:bg-amber-600 text-white p-2 rounded-full absolute right-10 top-10"
      >
        <Plus size={24} />
      </button>

      <div className="container w-[980px] min-h-[550px] mt-5 bg-white/30 rounded-[25px] shadow-xl p-6">
        {selected.length > 0 && (
          <div className="flex justify-end gap-4 mb-2 items-center">
            <div className=" text-black px-4 py-1 rounded-md">
              Selected: {selected.length} item{selected.length > 1 ? "s" : ""}
            </div>

            <button
              onClick={handleDelete}
              className=" bg-red-500 text-white px-1 py-1 rounded cursor-pointer hover:bg-red-800"
            >
              <Trash2Icon />
            </button>
          </div>
        )}

        <div className="flex flex-wrap gap-4 relative">
          {locations.length === 0 ? (
            <p className="text-gray-500 italic">There's no location at all.</p>
          ) : (
            locations.map((item) => (
              <div key={item.id} className="w-[300px] h-[240px] perspective">
                <div
                  className={`relative w-full h-full duration-700 transform-style-preserve-3d ${
                    flippedCard === item.id ? "rotate-y-180" : ""
                  }`}
                >
                  <div className="absolute w-full h-full backface-hidden bg-white/70 rounded-lg p-4 shadow-md">
                    <input
                      type="checkbox"
                      className="absolute top-1 left-1"
                      checked={selected.includes(item.id)}
                      onChange={() => toggleSelect(item.id)}
                    />
                    <img
                      src={getWeatherIconUrl(item.code)}
                      alt={getWeatherDescription(item.code)}
                      className="w-15 h-15 absolute right-5"
                    />
                    <h2 className="text-xl font-bold mt-1">{item.name}</h2>
                    <p className="text-sm text-gray-700 italic">
                      Lat: {item.lat}, Lon: {item.lon}
                    </p>
                    <p className="mt-2 text-gray-800">{item.summary}</p>
                    <p className="text-gray-800">
                      Temp:{" "}
                      {item.temperature !== null
                        ? `${item.temperature}°C`
                        : "-"}
                    </p>
                    <p className="text-gray-800">
                      Wind:{" "}
                      {item.windspeed !== null ? `${item.windspeed} km/h` : "-"}
                    </p>
                    <p className="text-xs text-gray-500 mt-1 italic">
                      Updated at: {item.updatedAt}
                    </p>
                    <div className="flex justify-between mt-4">
                      <button
                        onClick={() => {
                          setEditingLocation(item);
                          setShowModal(true);
                        }}
                        className="text-blue-600 hover:text-blue-800 flex items-center gap-1"
                      >
                        <Pencil size={16} /> Edit
                      </button>
                      <button
                        onClick={() => toggleFlip(item.id)}
                        className="text-amber-600 hover:text-amber-800 flex items-center gap-1"
                      >
                        <BarChart2 size={16} /> Prediction
                      </button>
                    </div>
                  </div>

                  <div className="absolute w-full h-full backface-hidden rotate-y-180 bg-white/70 rounded-lg p-4 shadow-md overflow-auto">
                    <h2 className="text-lg font-bold mb-2">Prediction</h2>
                    {item.predictions.length === 0 ? (
                      <p className="text-sm text-gray-700 italic">
                        No predictions available.
                      </p>
                    ) : (
                      item.predictions.map((pred: any, index: number) => (
                        <div key={index} className="mb-2">
                          <p className="text-sm font-medium">
                            {new Date(pred.Date).toLocaleString()}
                          </p>
                          <p className="text-sm">{pred.Summary}</p>
                        </div>
                      ))
                    )}
                    <button
                      onClick={() => toggleFlip(item.id)}
                      className="mt-4 text-blue-600 hover:text-blue-800"
                    >
                      Back
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {showModal && (
        <LocationModal
          mode={editingLocation ? "edit" : "add"}
          initialData={
            editingLocation
              ? {
                  name: editingLocation.name,
                  lat: editingLocation.lat,
                  lon: editingLocation.lon,
                }
              : undefined
          }
          onClose={() => {
            setShowModal(false);
            setEditingLocation(null);
          }}
          onSubmit={async (name, lat, lon) => {
            try {
              const token = localStorage.getItem("token");
              if (editingLocation) {
                await axios.put(
                  `http://localhost:8080/locations/${editingLocation.id}`,
                  { name, latitude: lat, longitude: lon },
                  {
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                    withCredentials: true,
                  }
                );
                toast.success("Lokasi berhasil diperbarui!");
              } else {
                await axios.post(
                  "http://localhost:8080/locations",
                  { name, latitude: lat, longitude: lon },
                  {
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                    withCredentials: true,
                  }
                );
                toast.success("Lokasi berhasil ditambahkan!");
              }
              setShowModal(false);
              setEditingLocation(null);
              setTimeout(() => window.location.reload(), 2000);
            } catch (error) {
              toast.error(
                editingLocation
                  ? "Gagal memperbarui lokasi."
                  : "Gagal menambahkan lokasi."
              );
              console.error(error);
            }
          }}
        />
      )}
    </div>
  );
}
