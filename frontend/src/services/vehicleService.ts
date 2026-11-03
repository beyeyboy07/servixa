
import api from "./api";

export interface Vehicle {
  id: number;
  plate_number: string;
  brand: string;
  model: string;
  year: number;
  status: string;
  is_active: boolean;
}

export interface CreateVehiclePayload {
  plate_number: string;
  brand: string;
  model: string;
  year: number;
}

export interface UpdateVehiclePayload {
  plate_number: string;
  brand: string;
  model: string;
  year: number;
  status: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

// Kendaraan yang tersedia untuk assignment Work Order
export const getAvailableVehicles = async (): Promise<Vehicle[]> => {
  const response = await api.get<ApiResponse<Vehicle[]>>("/vehicles");

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve available vehicles",
    );
  }

  return response.data.data;
};

// Seluruh kendaraan untuk halaman Fleet Management
export const getAllVehicles = async (): Promise<Vehicle[]> => {
  const response = await api.get<ApiResponse<Vehicle[]>>("/vehicles/all");

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve vehicles",
    );
  }

  return response.data.data;
};

export const createVehicle = async (
  payload: CreateVehiclePayload,
): Promise<Vehicle> => {
  const response = await api.post<ApiResponse<Vehicle>>("/vehicles", {
    plate_number: payload.plate_number.trim().toUpperCase(),
    brand: payload.brand.trim(),
    model: payload.model.trim(),
    year: payload.year,
  });

  if (!response.data.success) {
    throw new Error(response.data.message || "Failed to create vehicle");
  }

  return response.data.data;
};

export const updateVehicle = async (
  id: number,
  payload: UpdateVehiclePayload,
): Promise<Vehicle> => {
  const response = await api.put<ApiResponse<Vehicle>>(`/vehicles/${id}`, {
    plate_number: payload.plate_number.trim().toUpperCase(),
    brand: payload.brand.trim(),
    model: payload.model.trim(),
    year: payload.year,
    status: payload.status,
  });

  if (!response.data.success) {
    throw new Error(response.data.message || "Failed to update vehicle");
  }

  return response.data.data;
};

// Activate vehicle
export const activateVehicle = async (id: number): Promise<Vehicle> => {
  const response = await api.patch<ApiResponse<Vehicle>>(
    `/vehicles/${id}/activate`,
  );

  if (!response.data.success) {
    throw new Error(response.data.message || "Failed to activate vehicle");
  }

  return response.data.data;
};

// Deactivate vehicle
export const deactivateVehicle = async (id: number): Promise<Vehicle> => {
  const response = await api.patch<ApiResponse<Vehicle>>(
    `/vehicles/${id}/deactivate`,
  );

  if (!response.data.success) {
    throw new Error(response.data.message || "Failed to deactivate vehicle");
  }

  return response.data.data;
};