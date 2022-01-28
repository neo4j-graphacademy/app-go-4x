package neoflix;

import com.google.gson.*;

import java.lang.reflect.Type;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.List;

public class GsonUtils {
    public static Gson gson() throws ClassNotFoundException {
        GsonBuilder gsonBuilder = new GsonBuilder();
        gsonBuilder.registerTypeAdapter(LocalDate.class, new LocalDateSerializer());
        Class type = Class.forName("java.util.Collections$EmptyList");
        gsonBuilder.registerTypeAdapter(type, new EmptyListSerializer());
        return gsonBuilder.create();
    }

    static class LocalDateSerializer implements JsonSerializer<LocalDate> {
        private static final DateTimeFormatter formatter = DateTimeFormatter.ofPattern("dd-MMM-yyyy");

        @Override
        public JsonElement serialize(LocalDate localDate, Type srcType, JsonSerializationContext context) {
            return new JsonPrimitive(formatter.format(localDate));
        }
    }
    static class EmptyListSerializer implements JsonSerializer<List> {

        @Override
        public JsonElement serialize(List list, Type srcType, JsonSerializationContext context) {
            return new JsonArray(0);
        }
    }
}
