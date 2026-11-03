
import api from "./api";

export type ServiceRequestStatus =
  | "PENDING"
  | "REVIEWED"
  | "REJECTED";

export interface ServiceRequest {
  id: number;
  customer_id: number;
  title: string;
  description: string;
  service_address: string;
  preferred_date: string | null;
  status: ServiceRequestStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateServiceRequestPayload {
  title: string;
  description: string;
  service_address: string;
  preferred_date: string | null;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export const getMyServiceRequests = async (): Promise<ServiceRequest[]> => {
  const response = await api.get<ApiResponse<ServiceRequest[]>>(
    "/service-requests",
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve service requests",
    );
  }

  return response.data.data;
};

export const createServiceRequest = async (
  payload: CreateServiceRequestPayload,
): Promise<ServiceRequest> => {
  const response = await api.post<ApiResponse<ServiceRequest>>(
    "/service-requests",
    payload,
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to create service request",
    );
  }

  return response.data.data;
};

export const reviewServiceRequest = async (
  id: number,
): Promise<ServiceRequest> => {
  const response = await api.patch<ApiResponse<ServiceRequest>>(
    `/service-requests/${id}/review`,
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to review service request",
    );
  }

  return response.data.data;
};

export const rejectServiceRequest = async (
  id: number,
): Promise<ServiceRequest> => {
  const response = await api.patch<ApiResponse<ServiceRequest>>(
    `/service-requests/${id}/reject`,
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to reject service request",
    );
  }

  return response.data.data;
};
