using SimpleJson;


	public static class JsonDataExtension
	{
        public static string AsStr(this JsonObject self, string key, string defValue)
        {
            try
            {
                return self[key].ToString();
            }
            catch
            {
                return defValue;
            }
        }

        public static string AsStr(this JsonObject self, int index, string defValue)
        {
            try
            {
                return self[index].ToString();
            }
            catch
            {
                return defValue;
            }
        }
        
        public static bool AsBool(this JsonObject self, string key, bool defValue)
        {
            try
            {
                return int.Parse(self[key].ToString()) == 1;
            }
            catch
            {
                return defValue;
            }
        }

        public static bool AsBool(this JsonObject self, int index, bool defValue)
        {
            try
            {
                return int.Parse(self[index].ToString()) == 1;
            }
            catch
            {
                return defValue;
            }
        }

        public static int AsInt(this JsonObject self, string key, int defValue)
        {
            try
            {
                return int.Parse(self[key].ToString());
            }
            catch
            {
                return defValue;
            }
        }

        public static int AsInt(this JsonObject self, int index, int defValue)
        {
            try
            {
                return int.Parse(self[index].ToString());
            }
            catch
            {
                return defValue;
            }
        }

        public static float AsFloat(this JsonObject self, string key, float defValue)
        {
            try
            {
                return float.Parse(self[key].ToString());
            }
            catch
            {
                return defValue;
            }
        }

        public static float AsFloat(this JsonObject self, int index, float defValue)
        {
            try
            {
                return float.Parse(self[index].ToString());
            }
            catch
            {
                return defValue;
            }
        }

        public static JsonObject GetChild(this JsonObject self, string key)
        {
            try
            {
                return self[key] as JsonObject;
            }
            catch
            {
                return null;
            }
        }

        public static JsonObject GetChild(this JsonObject self, int index)
        {
            try
            {
                return self[index] as JsonObject;
            }
            catch
            {
                return null;
            }
        }
	}

