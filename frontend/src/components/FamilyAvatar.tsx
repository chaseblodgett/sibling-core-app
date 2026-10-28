import { Avatar } from "@mui/material";
import type { FamilyColor } from "../utils/familyColor";

interface FamilyAvatarProps {
  name: string;
  color: FamilyColor;
  size?: number;
}

export function FamilyAvatar({ name, color, size = 32 }: FamilyAvatarProps) {
  const initial = name.trim().charAt(0).toUpperCase() || "?";
  return (
    <Avatar
      sx={{
        bgcolor: color.border,
        color: "#fff",
        width: size,
        height: size,
        fontSize: Math.round(size * 0.45),
        fontWeight: 600,
      }}
    >
      {initial}
    </Avatar>
  );
}
