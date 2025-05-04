import axios from "axios";
import { useState } from "react";
import toast from "react-hot-toast";

function ScheduleModal({
  onClose,
  selectedIds,
}: {
  onClose: () => void;
  selectedIds: number[];
}) {
  const [duration, setDuration] = useState("12");
  const [unit, setUnit] = useState("h");

  const handleSubmit = async () => {
    try {
      const token = localStorage.getItem("token");
      await axios.post(
        `http://localhost:8080/schedule`,
        {
          ids: selectedIds,
          interval: `${duration}${unit}`, // misalnya: "12h"
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );
      toast.success("Schedule berhasil diperbarui.");
      onClose();
    } catch (err) {
      toast.error("Gagal memperbarui schedule.");
      console.error(err);
    }
  };

  return (
    <div className="modal">
      <h2>Set Schedule Interval</h2>
      <input
        type="number"
        value={duration}
        onChange={(e) => setDuration(e.target.value)}
        className="input"
      />
      <select
        value={unit}
        onChange={(e) => setUnit(e.target.value)}
        className="select"
      >
        <option value="s">Detik</option>
        <option value="m">Menit</option>
        <option value="h">Jam</option>
      </select>
      <button onClick={handleSubmit}>Submit</button>
      <button onClick={onClose}>Cancel</button>
    </div>
  );
}

export default ScheduleModal;
