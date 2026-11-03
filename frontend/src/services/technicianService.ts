import api from "./api";

export interface Technician {
  id: number;
  user_id: number;
  name: string;
  email: string;
  phone: string;
  specialization: string;
  is_available: boolean;
  is_active: boolean;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export const getAvailableTechnicians = async (): Promise<Technician[]> => {
  const response = await api.get<ApiResponse<Technician[]>>(
    "/technicians",
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve technicians",
    );
  }

  return response.data.data;
};